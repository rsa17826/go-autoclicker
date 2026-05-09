{
  description = "A simple Go autoclicker";

  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-unstable";
    };
    utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = {
          # The actual package
          default = pkgs.buildGoModule {
            pname = "go-autoclicker";
            version = "1";
            src = ./.;
            vendorHash = "sha256-ePH+0XpQLVVJG8eZv+Uhxjmd9Du0KYSBjOAboFM34M4=";
          };
        };
        devShells = {
          # Development environment
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
            ];
          };
        };
      }
    );
}
