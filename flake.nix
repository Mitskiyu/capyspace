{
  description = "capyspace";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    forAllSystems = function: nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed (system: function nixpkgs.legacyPackages.${system});
  in {
    devShells = forAllSystems (pkgs: {
      default = pkgs.mkShellNoCC {
        packages = with pkgs; [
          bun
          go
        ];
      };
    });
  };
}
