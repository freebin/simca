<?php
$options = getopt('c:n:');
if (empty($options['c']) || ($options['c'] < 0) || empty($options['n']) || ($options['n'] < 0)) {
    echo "Usage: test.php -c<CONCURRENCY> -n<ITERATIONS>\n";
    exit;
}

if (!extension_loaded('pcntl')) {
    dl('pcntl.so');
}

$job_chunks = [];
for ($i = 0; $i < $options['c']; $i++) {
    $chunk = [];
    for ($j = 0; $j < $options['n']; $j++) {
        $chunk[] = [
            'op' => rand(0, 3) ? 'get' : 'set', //TODO: Adjustable proportion
            'key' => "test_key_" . rand(0, $options['n']), //TODO: Controllable/adjustable conflicts
        ];
    }

    $job_chunks[$i] = $chunk;
}

$pids = [];
for ($i = 0; $i < $options['c']; $i++) {
    $pid = pcntl_fork();
    if ($pid == 0) {
        doZeeJob($job_chunks[$i]);
        exit;
    } elseif ($pid < 0) {
        echo "Failed to fork\n";
        exit;
    } else {
        $pids[] = $pid;
    }
}

echo "Waiting for kids\n";
foreach ($pids as $pid) {
    pcntl_waitpid($pid, $status);
    unset($pids[$pid]);
}
echo "All kids finished\n";

function doZeeJob($chunk)
{
    $socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
    if (!$socket) {
        echo "[FATAL] Failed to create socket" . PHP_EOL;
        exit -1;
    }

    $host = 'localhost';
    $port = 8080;

    echo 'Connecting to ' . $host . ':' . $port . PHP_EOL;
    $res = socket_connect($socket, $host, $port);
    if (!$res) {
        echo "[FATAL] Failed to connect to demon" . PHP_EOL;
        exit -1;
    }

    foreach ($chunk as $command) {
        $command_str = $command['op'] . ' ' . $command['key'] . ($command['op'] == 'set' ? ' ' . rand(1000, 9999) : '');

        $res = socket_write($socket, $command_str, strlen($command_str));
        if (!$res) {
            echo "[FATAL] Write to demon failed while sending '" . $command_str . "'" . PHP_EOL;
            exit -1;
        }

        $resp = '';
        while ($out = socket_read($socket, 2048, PHP_NORMAL_READ)) {
            $resp .= $out;
            if (strpos($out, "\n") !== false) {
                break;
            }
        }

        //echo "[RESP] " . $resp . "\n";
    }

    $res = socket_write($socket, "exit", strlen($command_str));
    socket_close($socket);
}