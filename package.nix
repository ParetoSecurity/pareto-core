{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-TlZ7Z6qZCdmXej9oaB4ImnVuP1AKoLhDIqM0ga1V/O8=";
  subPackages = ["cmd/paretosecurity"];
  doCheck = true;
}
