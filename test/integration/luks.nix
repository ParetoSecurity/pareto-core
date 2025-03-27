let
  # Override paretosecurity to use the local codebase
  pareto = {
    pkgs,
    lib,
    ...
  }: {
    services.paretosecurity = {
      enable = true;
      package = pkgs.callPackage ../../package.nix {inherit lib;};
    };
  };

  # Easier tests debugging by SSH-ing into nodes
  ssh = {port}: {...}: {
    services.openssh = {
      enable = true;
      settings = {
        PermitRootLogin = "yes";
        PermitEmptyPasswords = "yes";
      };
    };
    security.pam.services.sshd.allowNullPassword = true;
    virtualisation.forwardPorts = [
      {
        from = "host";
        host.port = port;
        guest.port = 22;
      }
    ];
  };
in {
  name = "Check test: Filesystem encryption is enabled";

  nodes = {
    plaindisk = {
      pkgs,
      lib,
      ...
    }: {
      imports = [
        (pareto {inherit pkgs lib;})
      ];
    };

    luks = {
      pkgs,
      lib,
      config,
      ...
    }: {
      imports = [
        (pareto {inherit pkgs lib;})
      ];

      # NixOS LUKS test VM configuration taken from
      # https://github.com/NixOS/nixpkgs/blob/master/nixos/tests/luks.nix

      boot.initrd.extraUtilsCommands = ''
        # We need mke2fs in the initrd.
        copy_bin_and_libs ${pkgs.e2fsprogs}/bin/mke2fs
      '';

      boot.initrd.postDeviceCommands = ''
        # If the disk image appears to be empty, run mke2fs to
        # initialise.
        FSTYPE=$(blkid -o value -s TYPE ${config.virtualisation.rootDevice} || true)
        PARTTYPE=$(blkid -o value -s PTTYPE ${config.virtualisation.rootDevice} || true)
        if test -z "$FSTYPE" -a -z "$PARTTYPE"; then
            mke2fs -t ext4 ${config.virtualisation.rootDevice}
        fi
      '';

      # Use systemd-boot
      virtualisation = {
        emptyDiskImages = [
          512
          512
        ];
        useBootLoader = true;
        useEFIBoot = true;
        # To boot off the encrypted disk, we need to have a init script which comes from the Nix store
        mountHostNixStore = true;
      };
      boot.loader.systemd-boot.enable = true;

      boot.kernelParams = lib.mkOverride 5 ["console=tty1"];

      environment.systemPackages = with pkgs; [cryptsetup];

      specialisation = rec {
        boot-luks.configuration = {
          boot.initrd.luks.devices = lib.mkVMOverride {
            # We have two disks and only type one password - key reuse is in place
            cryptroot.device = "/dev/vdb";
            cryptroot2.device = "/dev/vdc";
          };
          virtualisation.rootDevice = "/dev/mapper/cryptroot";
        };
        boot-luks-custom-keymap.configuration = lib.mkMerge [
          boot-luks.configuration
          {
            console.keyMap = "neo";
          }
        ];
      };
    };
  };

  interactive.nodes.regularboot = {...}:
    ssh {port = 2221;} {};

  interactive.nodes.secureboot = {...}:
    ssh {port = 2222;} {};

  enableOCR = true;

  testScript = {nodes, ...}: ''
    # Test 1: check fails with LUKS not configured
    out = plaindisk.fail("paretosecurity check --only 21830a4e-84f1-48fe-9c5b-beab436b2cdb")
    expected = (
        "  • Starting checks...\n"
        "  • System Integrity: Filesystem encryption is enabled > [FAIL] Block device encryption is disabled\n"
        "  • Checks completed.\n"
    )
    assert out == expected, f"Expected did not match actual, got \n{out}"

    # Test 2: check succeeds with LUKS configured
    luks.wait_for_unit("multi-user.target")
    luks.succeed("echo -n supersecret | cryptsetup luksFormat -q --iter-time=1 /dev/vdb -")
    luks.succeed("echo -n supersecret | cryptsetup luksFormat -q --iter-time=1 /dev/vdc -")
    luks.succeed("bootctl set-default nixos-generation-1-specialisation-boot-luks.conf")
    luks.succeed("sync")
    luks.crash()
    luks.start()
    luks.wait_for_text("Passphrase for")
    luks.send_chars("supersecret\n")
    luks.wait_for_unit("multi-user.target")
    assert "/dev/mapper/cryptroot on / type ext4" in luks.succeed("mount")

    out = luks.succeed("paretosecurity check --only 21830a4e-84f1-48fe-9c5b-beab436b2cdb")
    expected = (
        "  • Starting checks...\n"
        "  • System Integrity: Filesystem encryption is enabled > [OK] \n"
        "  • Checks completed.\n"
    )
    assert out == expected, f"Expected did not match actual, got \n{out}"
  '';
}
