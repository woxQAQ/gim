let
  pkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/6607cf789e541e7873d40d3a8f7815ea92204f32.tar.gz") {};
in
  pkgs.mkShell {
    buildInputs = with pkgs; [go_1_23];
    hardeningDisable = ["fortify"];
  }
