assert "Pareto Security CLI" in machine.succeed("paretosecurity --help")

res = machine.fail("paretosecurity check")
fail_count = res.count("✗")
dial_error_count = res.count("Failed to connect to helper")
assert (
    dial_error_count == 0
), f"Helper could not start, found : {dial_error_count} calls to dial error"
assert fail_count > 3, f"Found {fail_count} failed checks"
