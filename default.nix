{
  buildGoModule,
  fetchFromGitHub,
  installShellFiles,
  stdenv,
  lib,
}:

buildGoModule rec {
  pname = "talhelper";
  version = "2.4.3";

  src = fetchFromGitHub {
    owner = "budimanjojo";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-NOU8A/6gX5JINu2+cSbsg/1KAX8B8hWfjf9eb/Azbm4=";
  };

  vendorHash = "sha256-xapCHIquqH1WrYc7ZJY1jjV5D5eaeO1bNirMxO0d+EA=";

  ldflags = [ "-s -w -X github.com/budimanjojo/talhelper/cmd.version=v${version}" ];

  doCheck = false; # no tests

  subPackages = [
    "."
    "./cmd"
  ];

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
