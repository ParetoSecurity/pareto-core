{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-APxHPQ+08flAWdgVLi7BAJ4D/8klWzKIF+v7dlC3y28=";
  doCheck = true;
}
