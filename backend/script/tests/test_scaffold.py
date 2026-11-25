import importlib.util
import json
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

SCRIPT = Path(__file__).resolve().parents[1] / "new_service.py"
SPEC = importlib.util.spec_from_file_location("new_service", SCRIPT)
scaffold = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(scaffold)


class ScaffoldTests(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.output = Path(temporary.name) / "service"

    def args(self, *extra):
        return ["--name", "test-service", "--module", "example.com/test-service",
                "--out", str(self.output), *extra]

    def run_generator(self, *extra):
        return subprocess.run([sys.executable, str(SCRIPT), *self.args(*extra)],
                              capture_output=True, text=True)

    def test_invalid_name_and_port(self):
        for args in [("--name", "../other"), ("--port", "65536"), ("--port", "abc")]:
            with self.subTest(args=args):
                self.assertNotEqual(self.run_generator(*args).returncode, 0)
                self.assertFalse(self.output.exists())

    def test_force_refuses_unmarked_directory(self):
        self.output.mkdir()
        original = self.output / "original"
        original.write_text("keep")
        self.assertNotEqual(self.run_generator("--force").returncode, 0)
        self.assertEqual(original.read_text(), "keep")

    def test_collision_requires_force(self):
        self.assertEqual(self.run_generator().returncode, 0)
        self.assertNotEqual(self.run_generator().returncode, 0)
        self.assertEqual(self.run_generator("--force").returncode, 0)

    def test_render_failure_preserves_previous_scaffold(self):
        self.assertEqual(self.run_generator().returncode, 0)
        original = self.output / "README.md"
        original.write_text("keep")
        with patch.object(sys, "argv", [str(SCRIPT), *self.args("--force")]):
            with patch.object(scaffold, "process_file", side_effect=OSError("render failed")):
                with self.assertRaises(SystemExit):
                    scaffold.main()
        self.assertEqual(original.read_text(), "keep")

    def test_publish_failure_restores_previous_scaffold(self):
        self.assertEqual(self.run_generator().returncode, 0)
        original = self.output / "README.md"
        original.write_text("keep")
        rename = Path.rename

        def fail_on_publish(source, target):
            if target == self.output.resolve() and source.parent.name.startswith(".scaffold-backup-"):
                return rename(source, target)
            if target == self.output.resolve():
                raise OSError("publish failed")
            return rename(source, target)

        with patch.object(sys, "argv", [str(SCRIPT), *self.args("--force")]):
            with patch.object(Path, "rename", fail_on_publish):
                with self.assertRaisesRegex(SystemExit, "publish failed"):
                    scaffold.main()
        self.assertEqual(original.read_text(), "keep")

    def test_force_refuses_symlink_and_git_root(self):
        self.assertEqual(self.run_generator().returncode, 0)
        link = self.output.with_name("link")
        link.symlink_to(self.output, target_is_directory=True)
        result = self.run_generator("--force", "--out", str(link))
        self.assertNotEqual(result.returncode, 0)
        self.assertTrue(link.is_symlink())
        (self.output / ".git").mkdir()
        self.assertNotEqual(self.run_generator("--force").returncode, 0)

    @unittest.skipUnless(shutil.which("docker"), "Docker Compose is required")
    def test_compose_variants(self):
        variants = {
            "minimal": ([], {"test-service"}),
            "redis": (["--cache", "redis"], {"test-service", "redis"}),
            "full": (
                ["--port", "9191", "--cache", "valkey", "--db", "postgres",
                 "--kafka", "yes", "--storage", "s3", "--firebase", "yes",
                 "--mockoon", "yes", "--mocktcp", "yes"],
                {"test-service", "redis", "postgres", "kafka", "kafka-ui", "minio", "mockoon"},
            ),
        }
        for name, (flags, expected) in variants.items():
            with self.subTest(name=name):
                output = self.output.with_name(name)
                generated = self.run_generator("--out", str(output), *flags)
                self.assertEqual(generated.returncode, 0, generated.stderr)
                subprocess.run(["docker", "compose", "config", "--quiet"], cwd=output, check=True)
                config = subprocess.check_output(
                    ["docker", "compose", "config", "--format", "json"], cwd=output, text=True,
                )
                services = json.loads(config)["services"]
                self.assertEqual(set(services), expected)
                if name == "redis":
                    self.assertEqual(services["redis"]["image"], "redis:7-alpine")
                if name == "full":
                    cache = services["redis"]
                    self.assertTrue(cache["image"].startswith("valkey/valkey:"))
                    self.assertEqual(cache["command"][0], "valkey-server")
                    self.assertIn("valkey-cli", cache["healthcheck"]["test"])
                    self.assertIn("service-network", services["minio"]["networks"])
                    mockoon = services["mockoon"]
                    mount = mockoon["volumes"][0]
                    dataset = mockoon["command"][1]
                    self.assertTrue((Path(mount["source"]) / Path(dataset).name).is_file())

    def test_custom_port_and_name(self):
        result = self.run_generator("--port", "9191", "--kafka", "yes")
        self.assertEqual(result.returncode, 0, result.stderr)
        for file, expected in {
            "Dockerfile": ["EXPOSE 9191", "localhost:9191", "./test-service"],
            "Makefile": ["SERVICE_NAME ?= test-service", "-p 9191:9191"],
            "docker-compose.yaml": ['"9191:9191"', '${KAFKA_UI_PORT:-8090}:8080'],
            "config/config.docker-compose.yaml": ['port: "9191"'],
            "bruno/environments/LOCAL.yml": ["localhost:9191"],
        }.items():
            content = (self.output / file).read_text()
            self.assertNotIn("TODO-service-port", content)
            for value in expected:
                self.assertIn(value, content)


if __name__ == "__main__":
    unittest.main()
