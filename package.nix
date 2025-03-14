{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-XK5ma3c84OakVYEVkEropiNU1KaZsdPpiYw0dkTdfdY=";
  subPackages = ["cmd/paretosecurity"];
  doCheck = true;
}
