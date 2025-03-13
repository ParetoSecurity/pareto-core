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

  # Override build step with custom ldflags
  buildPhase = ''
    runHook preBuild

    buildGoDir() {
      local dir="$1"
      cd "$dir"
      go build -ldflags '
        -s -w
        -X github.com/ParetoSecurity/agent/shared.Version=${version}
        -X github.com/ParetoSecurity/agent/shared.Commit=${builtins.substring 0 7 (lib.commitIdFromGitRepo ./.git)}
        -X github.com/ParetoSecurity/agent/shared.Date=${builtins.substring 0 10 builtins.currentTime}
      ' -v -o $GOPATH/bin/$(basename $dir) .
      cd -
    }

    buildGoDir ./cmd/paretosecurity

    runHook postBuild
  '';
}
