<?hh

class Config {
  public string $server_address;
  public int $max_attempt;
  public int $timeout;
  public string $json_str;

  public function __construct(string $confFilename): void {
    $this->json_str = file_get_contents($confFilename);
    $json = json_decode($this->json_str, true);
    $this->server_address = $json['address'];
    $this->max_attempt = $json['max_attempts'];
    $this->timeout = $json['timeout'];
  }

  public function printStats(): void {
    echo $this->json_str;
  }
}

class Client {

  private Config $config;
  private bool $waiting_for_input = false;

  public function printUsage(): void {
    echo "b                <Number of tickets>\n";
    echo "q/quit/e/exit    quit\n";

  }

  public function __construct(string $confFilename) {
    $this->config = new Config($confFilename);
  }

  public function run(): void {
    $this->printUsage();
    for (; ; ) {
      $this->waiting_for_input = true;
      $cmd = readline("> ");
      dispatch($cmd);
    }
  }

  private function dispatch(string $cmd): void {
    $tokens = preg_split("/\s+/", $cmd);

    var_dump($tokens);
    switch ($favcolor) {
      case "red":
        echo "Your favorite color is red!";
        break;
      case "blue":
        echo "Your favorite color is blue!";
        break;
      case "green":
        echo "Your favorite color is green!";
        break;
      default:
        echo "Your favorite color is neither red, blue, nor green!";
    }

    if ($tokens[0] == 'q') {
      return;
    }
  }
}

$client = new Client("conf.json");
$client->run();
