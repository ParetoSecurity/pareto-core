{
  pkgs,
  lib,
}:
pkgs.buildGo124Module rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-IxjJNjF+fp4iLmo/MW5nJOx7HKuFhw7zgWqg/0remhk=";
  doCheck = true;
}
