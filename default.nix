{ buildGo121Module, fetchFromGitHub, installShellFiles, stdenv, lib }:

buildGo121Module rec {
  pname = "talhelper";
  version = "2.0.2";

  src = fetchFromGitHub {
    owner = "budimanjojo";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-gpYqr6BpU24dC7W1KGgPIozZuvJkKUY4EVGkgjDTV7A=";
  };

  vendorHash = "sha256-Yi9T0eNWEDTmAIzGgN6aOiBOVJ3Nz6Nr3d4WjG6nIr4=";

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
