{
  flakePackage,
  pkgs ?
    import <nixpkgs> {
      inherit system;
      config = {};
    },
  system ? builtins.currentSystem,
}:
pkgs.testers.runNixOSTest {
  name = "pareto";
  nodes.machine = {pkgs, ...}: {
    environment.systemPackages = [flakePackage];

    systemd.sockets."paretosecurity" = {
      wantedBy = ["sockets.target"];
      socketConfig = {
        ListenStream = "/var/run/paretosecurity.sock";
        SocketMode = "0666";
      };
    };

    systemd.services."paretosecurity" = {
      requires = ["paretosecurity.socket"];
      after = ["paretosecurity.socket"];
      wantedBy = ["multi-user.target"];
      serviceConfig = {
        ExecStart = ["${flakePackage}/bin/paretosecurity" "helper" "--verbose" "--socket" "/var/run/paretosecurity.sock"];
        User = "root";
        Group = "root";
        StandardInput = "socket";
        Type = "oneshot";
        RemainAfterExit = "no";
        StartLimitInterval = "1s";
        StartLimitBurst = 100;
        ProtectSystem = "full";
        ProtectHome = true;
        StandardOutput = "journal";
        StandardError = "journal";
      };
    };

    systemd.user.services."pareto-core-hourly" = {
      wantedBy = ["timers.target"];
      serviceConfig = {
        Type = "oneshot";
        ExecStart = "${flakePackage}/bin/paretosecurity check";
        StandardInput = "null";
      };
    };

    systemd.user.timers."pareto-core-hourly" = {
      wantedBy = ["timers.target"];
      timerConfig = {
        OnCalendar = "hourly";
        Persistent = true;
      };
    };
  };

  interactive.nodes.machine = {...}: {
    services.openssh.enable = true;
    services.openssh.settings = {
      PermitRootLogin = "yes";
      PermitEmptyPasswords = "yes";
    };
    security.pam.services.sshd.allowNullPassword = true;
    virtualisation.forwardPorts = [
      {
        from = "host";
        host.port = 2222;
        guest.port = 22;
      }
    ];
  };

  testScript = builtins.readFile "${toString ./.}/nixos.py";
}
