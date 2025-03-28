{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-ZaJkRHGNw2azagJxmbn/EbBXGHGLTmfCU6toAmbHtM8=";
  doCheck = true;
}
