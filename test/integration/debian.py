vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed(
    "DEBIAN_FRONTEND=noninteractive sudo dpkg -i /mnt/package/paretosecurity_amd64.deb"
)

res = vm.succeed("paretosecurity check --json")
fail_count = res.count("fail")
dial_error_count = res.count("Failed to connect to helper")
assert (
    dial_error_count == 0
), f"Helper could not start, found : {dial_error_count} calls to dial error"
assert fail_count > 3, f"Found {fail_count} failed checks"
