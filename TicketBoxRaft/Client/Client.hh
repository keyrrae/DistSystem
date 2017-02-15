<?hh

class Config {
  public string $serverAddress;
  public int $maxAttempt;
  public int $timeout;

  public function __construct(string $confFilename): void {
    $jsonStr = file_get_contents($confFilename);
    $json = json_decode($jsonStr, true);
    $this->serverAddress = $json['address'];
    $this->maxAttempt = $json['max_attempts'];
    $this->timeout = $json['timeout'];

  }
}

class Client {

  private Config $config;

  public function printUsage(): void {
    echo "b <Number of tickets>\n";
  }

  public function __construct(string $confFilename) {
    $this->config = new Config($confFilename);
  }

  public function run(): void {
    $this->printUsage();
  }
}

$client = new Client("conf.json");
$client->run();
