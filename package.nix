{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-uyHSDtP4ChDfuB93yHqR+v4i2dnaCVawcNRb4ygfjUw=";
  subPackages = ["cmd/paretosecurity"];
  doCheck = true;
}
