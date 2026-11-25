# mocktcp fixtures

One subfolder per operation code your protocol defines, e.g.:

```
fixtures/
  5514/
    response.txt          # same response every call
  1801/
    response_1.txt        # round-robin: 1st call gets this
    response_2.txt        # 2nd call gets this
    response_3.txt        # 3rd call, then cycles back to response_1
```

Each `response*.txt` is the exact raw bytes to write back to the client —
capture real responses from the actual downstream (or its existing
mock, e.g. as400-adapter's `mockas400/fixtures/`) rather than
hand-constructing frames, since field offsets/padding are easy to get
subtly wrong by hand.

Delete this whole `mocktcp/` directory if your service's checklist answer
for "non-HTTP downstream" is no.
