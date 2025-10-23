{
  description = "bind_player flake";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "bind_player";
        version = "1.0.0";
        src = ./.;
        vendorHash = null;
        doCheck = false;
        subPackages = [ "." ];

        buildInputs = [
          pkgs.pkgconf
          pkgs.gcc
          pkgs.xorg.libX11
          pkgs.xorg.libXrandr
          pkgs.xorg.libXinerama
          pkgs.xorg.libXcursor
          pkgs.xorg.libXi
          pkgs.mesa
          pkgs.libGL
          pkgs.xorg.libXxf86vm 
        ];
      };
    };
}

