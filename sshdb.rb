class Sshdb < Formula
  homepage "https://github.com/trntv/sshdb"
  url "https://github.com/trntv/sshdb.git", :tag => "0.3.0"
  version "0.3.0"

  depends_on "go" => :build
  depends_on "dep" => :build

  def install
    ENV["GOPATH"] = buildpath
    (buildpath/"src/github.com/trntv/sshdb").install buildpath.children
    cd "src/github.com/trntv/sshdb" do
      system "dep ensure -vendor-only"
      system "make"
      bin.install Dir[buildpath/"bin/sshdb"]
    end
  end

end
