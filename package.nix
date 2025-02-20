{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-bCbykrAUqAhandPNAABbGfMewIPU1ojxVlvPenYtK38=";
  subPackages = ["cmd/paretosecurity"];
  doCheck = true;
}
