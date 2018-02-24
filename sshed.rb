class Sshed < Formula
  desc "SSH config and connections manager"
  homepage "https://github.com/trntv/sshed"
  url "https://github.com/trntv/sshed.git", :tag => "1.0.3"

  version "1.0.3"

  depends_on "go" => :build
  depends_on "dep" => :build

  def install
    ENV["GOPATH"] = buildpath
    (buildpath/"src/github.com/trntv/sshed").install buildpath.children
    cd "src/github.com/trntv/sshed" do
      system "dep ensure -vendor-only"
      system "make"
      bin.install Dir["build/sshed"]
      prefix.install Dir["completions/*"]
    end
  end

end
