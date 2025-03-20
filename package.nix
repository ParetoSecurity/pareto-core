{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-5DNwYEsfs0qGGwIUAH/zqxHgsa4MzzDEs9oP4zH/P6g=";
  doCheck = true;
}
