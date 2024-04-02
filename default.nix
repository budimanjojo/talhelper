{ buildGo122Module, fetchFromGitHub, installShellFiles, stdenv, lib }:

buildGo122Module rec {
  pname = "talhelper";
  version = "2.3.10";

  src = fetchFromGitHub {
    owner = "budimanjojo";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-pSOVleMZ+TwdJ0bxrF9FXI6jtC3HU/Eph9Nv3LD0LHc=";
  };

  vendorHash = "sha256-qWKWjglZgssEqZEBe5A15sSaiLNLCNNb7svjeGPhiIA=";

  ldflags = [ "-s -w -X github.com/budimanjojo/talhelper/cmd.version=v${version}" ];

  doCheck = false; # no tests

  subPackages = [ "." "./cmd" ];

  nativeBuildInputs = [ installShellFiles ];

  doInstallCheck = true;
  installCheckPhase = ''
    $out/bin/talhelper --version | grep ${version} > /dev/null
  '';

  postInstall = lib.optionalString (stdenv.hostPlatform == stdenv.buildPlatform) ''
    for shell in bash fish zsh; do
      $out/bin/talhelper completion $shell > talhelper.$shell
      installShellCompletion talhelper.$shell
    done
  '';

  meta = with lib; {
    description = "A tool to help creating Talos kubernetes cluster";
    longDescription = ''
      Talhelper is a helper tool to help creating Talos Linux cluster
      in your GitOps repository.
    '';
    homepage = "https://github.com/budimanjojo/talhelper";
    license = licenses.bsd3;
  };
}
