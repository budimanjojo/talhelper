{ buildGoModule, fetchFromGitHub, installShellFiles, stdenv, lib }:

buildGoModule rec {
  pname = "talhelper";
  version = "1.6.4";

  src = fetchFromGitHub {
    owner = "budimanjojo";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-VBqdRWKuny0oVKjfcJPrfP1SHQJM58+EDKwfQtEK9y0=";
  };

  vendorSha256 = "sha256-pPzkJcEbR3DXJo26N1fxwWq6tUmNQINH97J3rqvwTH4=";

  ldflags = [ "-s -w -X github.com/budimanjojo/talhelper/cmd.version=v${version}" ];

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
