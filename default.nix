{
  buildGo125Module,
  fetchFromGitHub,
  installShellFiles,
  stdenv,
  lib,
}:

buildGo125Module rec {
  pname = "talhelper";
  version = "3.1.4";

  src = fetchFromGitHub {
    owner = "budimanjojo";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-UxPWSd6lIpz754ygGRrLxhgS86qHt6KJ3sdCWWDF6Bw=";
  };

  vendorHash = "sha256-idtOLMg0KFUXTDl+y3uDYZekHMIyEVOzyZyNE67vXgY=";

  ldflags = [ "-s -w -X github.com/budimanjojo/talhelper/v3/cmd.version=v${version}" ];

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
    mainProgram = "talhelper";
    longDescription = ''
      Talhelper is a helper tool to help creating Talos Linux cluster
      in your GitOps repository.
    '';
    homepage = "https://github.com/budimanjojo/talhelper";
    license = licenses.bsd3;
  };
}
