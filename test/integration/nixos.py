
help = machine.succeed("paretosecurity --help")
assert "Pareto Security CLI" in help
print(help)

status, res = machine.execute("paretosecurity check")
assert "Failed to connect to helper" not in res, "Helper could not start"

res = machine.succeed("paretosecurity info")
fail_count = res.count("State: false")
print(res)


assert fail_count > 0, "End line not found"
