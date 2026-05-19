#!/usr/bin/env python3
"""Interop check: pylogix client against the gologix Server_Class3 example.

Run the gologix server first (`go run .` in this directory), then point this
script at it on the host where the server is listening. It exercises the two
STRING paths the issue #58 fix touches:

  1. Read `teststring`, which the server seeds with "Hello World", and
     assert the value comes back intact.
  2. Write a fresh value to `writestring`, read it back, and assert it
     round-trips.

The script does not catch every CIP edge case — its job is to confirm that a
non-gologix client (pylogix) accepts the server's STRING wire encoding for
reads and that the server parses pylogix's STRING write payload correctly.

Requires pylogix: `pip install pylogix`.
"""

import argparse
import sys


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--host",
        default="127.0.0.1",
        help="IP address of the gologix Server_Class3 process (default: 127.0.0.1).",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=44818,
        help="EIP port (default 44818). Override when using a capture proxy.",
    )
    parser.add_argument(
        "--write-value",
        default="pylogix-round-trip-check",
        help="Value to write to the writestring tag.",
    )
    args = parser.parse_args()

    try:
        import pylogix
    except ImportError:
        print("pylogix is not installed. Run: pip install pylogix", file=sys.stderr)
        return 2

    failures = 0
    with pylogix.PLC(args.host, port=args.port) as plc:
        # 1. Read the seeded STRING.
        resp = plc.Read("teststring")
        if resp.Status != "Success":
            print(f"read teststring FAILED: status={resp.Status}", file=sys.stderr)
            failures += 1
        elif resp.Value != "Hello World":
            print(
                f"read teststring MISMATCH: want 'Hello World', got {resp.Value!r}",
                file=sys.stderr,
            )
            failures += 1
        else:
            print(f"read teststring OK: {resp.Value!r}")

        # 2. Write + read-back round trip.
        wresp = plc.Write("writestring", args.write_value)
        if wresp.Status != "Success":
            print(
                f"write writestring FAILED: status={wresp.Status}",
                file=sys.stderr,
            )
            failures += 1
        else:
            rresp = plc.Read("writestring")
            if rresp.Status != "Success":
                print(
                    f"read writestring after write FAILED: status={rresp.Status}",
                    file=sys.stderr,
                )
                failures += 1
            elif rresp.Value != args.write_value:
                print(
                    f"writestring round-trip MISMATCH: want {args.write_value!r}, got {rresp.Value!r}",
                    file=sys.stderr,
                )
                failures += 1
            else:
                print(f"writestring round-trip OK: {rresp.Value!r}")

    return 0 if failures == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
