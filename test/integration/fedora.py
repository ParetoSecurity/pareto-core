vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo dnf install -y /mnt/package/*x86_64.rpm")

res = vm.succeed("paretosecurity check")
fail_count = res.count("fail")
dial_error_count = res.count("Failed to connect to helper")
assert (
    dial_error_count == 0
), f"Helper could not start, found : {dial_error_count} calls to dial error"
assert fail_count > 3, f"Found {fail_count} failed checks"
