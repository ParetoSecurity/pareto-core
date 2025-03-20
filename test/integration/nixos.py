
help = machine.succeed("paretosecurity --help")
assert "Pareto Security CLI" in help
print(help)

status, res = machine.execute("paretosecurity check")
print(res)

assert "Failed to connect to helper" in res , f"Helper could not start"
assert "Checks completed" in res, f"End line not found"
