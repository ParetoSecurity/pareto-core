
help = machine.succeed("paretosecurity --help")
assert "Pareto Security CLI" in help
print(help)

status, res = machine.execute("paretosecurity check")
assert "Failed to connect to helper" not in res, "Helper could not start"

res = machine.succeed("paretosecurity info 2>&1")
fail_count = res.count("false")
print(res)

assert fail_count > 0, f"Failed to find any failed checks, found: {fail_count}"
