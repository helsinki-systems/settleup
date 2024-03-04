with import <nixpkgs> {}; let
  unstable = pkgs.unstable or (import <unstable> {});
in mkShell {
  name = "settleup";

  nativeBuildInputs = [
    gnumake

    go
    unstable.golangci-lint

    air
  ];

  shellHook = ''
    if [[ -n "$IN_LORRI_SHELL" ]]; then
      export GOPATH="$(dirname $IN_LORRI_SHELL)/.gopath"
    else
      export GOPATH="$PWD/.gopath"
    fi
  '';
}
