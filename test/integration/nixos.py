assert "Pareto Security CLI" in machine.succeed("paretosecurity --help")

res = machine.fail("paretosecurity check")
end_line = res.count("Checks completed")
dial_error_count = res.count("Failed to connect to helper")
assert (
    dial_error_count == 0
), f"Helper could not start, found : {dial_error_count} calls to dial error"
assert end_line > 0, f"End line not found, found : {res}"
