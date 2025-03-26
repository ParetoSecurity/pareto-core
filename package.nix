{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-JuUI450GTL6mbc/k9JsXg4D42dXYqNdZUWtbfJgw2mQ=";
  doCheck = true;
}
