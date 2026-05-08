{
  description = "A simple Go autoclicker";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
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
        # The actual package
        packages.default = pkgs.buildGoModule {
          pname = "go-autoclicker";
          version = "1";
          src = ./.;
          deleteVendor = true;
          vendorHash = "sha256-Z1lODTODCkWOa98uf6/s9/urDoNjWqDrAIVnCZstkbz=";
        };

        # Development environment
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];
        };
      }
    );
}
