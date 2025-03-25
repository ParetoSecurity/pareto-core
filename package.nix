{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-JJPfCJtbgcoobzT9zpiYpxKou+xXlU4EhYcQbPjpdW8=";
  doCheck = true;
}
