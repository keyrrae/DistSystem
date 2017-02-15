<?hh

class Client {
  public function printUsage(): void {
    echo "b <Number of tickets\n";
  }

  public function __construct(private string $conf) {}

  public function run(): void {
    $this->printUsage();
  }
}

$client = new Client("conf.json");
$client->run();
