{
  description = "Ellie Daily Planner CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devshell.url = "github:numtide/devshell";
    devshell.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, devshell }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f system);
    in {
      packages = forAllSystems (system:
        let pkgs = nixpkgs.legacyPackages.${system};
        in {
          default = pkgs.buildGoModule {
            pname = "ellie";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-JFrzxduL0Wr3+CGfAmJbAcaCWRP/vLF6nQWds2aamtw=";
          };
        }
      );

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          devshellPkgs = devshell.legacyPackages.${system};
        in {
          default = devshellPkgs.mkShell {
            name = "ellie-cli";
            packages = with pkgs; [ go gopls gotools go-tools ];
            commands = [
              { name = "build"; command = "go build -o ellie ./cmd/ellie"; help = "Build the CLI"; }
              { name = "test"; command = "go test ./..."; help = "Run tests"; }
              { name = "lint"; command = "staticcheck ./..."; help = "Run linter"; }
              { name = "format"; command = "gofmt -w ."; help = "Format code"; }
            ];
          };
        }
      );
    };
}
