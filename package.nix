{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-viaNG0RpUdXXe+ZNg69Y+u+rbzU44WbcbDOqMjCPqPc=";
  doCheck = true;
}
