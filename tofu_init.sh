#!/bin/sh

# 创建opentofu模版文件夹
mkdir tofu-emplate
# 进入
cd tofu-emplate

# 阿里云:
# -----------------------------------
# 创建阿里云opentofu模版文件夹
mkdir -p aliyun
# 进入
cd aliyun
# aliyun.tf文件并写入内容
cat << EOF > aliyun.tf
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
      version = "1.189.0"
    }
  }
}

provider "alicloud" {
  # Configuration options
}
EOF
# 退出
cd ..
echo "Folder created and aliyun.tf file written successfully."


# 腾讯云:
# -----------------------------------
# 创建腾讯云opentofu模版文件夹
mkdir -p tencent
# 进入
cd tencent
# tencentcloudstack.tf文件并写入内容
cat << EOF > tencentcloudstack.tf
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
      version = "1.81.16"
    }
  }
}

provider "tencentcloud" {
  # Configuration options
}
EOF
# 退出
cd ..
echo "Folder created and tencentcloudstack.tf file written successfully."


# 百度云:
# -----------------------------------
# 创建百度云opentofu模版文件夹
mkdir -p baidu
# 进入
cd baidu
# baidubce.tf文件并写入内容
cat << EOF > baidubce.tf
terraform {
  required_providers {
    baiducloud = {
      source = "baidubce/baiducloud"
      version = "1.19.35"
    }
  }
}

provider "baiducloud" {
  # Configuration options
}
EOF
# 退出
cd ..
echo "Folder created and baidubce.tf file written successfully."


# 华为云:
# -----------------------------------
# 创建华为云opentofu模版文件夹
mkdir -p huawei
# 进入
cd huawei
# huaweicloud.tf文件并写入内容
cat << EOF > huaweicloud.tf
terraform {
  required_providers {
    huaweicloud = {
      source = "huaweicloud/huaweicloud"
      version = "1.44.2"
    }
  }
}

provider "huaweicloud" {
  # Configuration options
}
EOF
# 退出
cd ..
echo "Folder created and huaweicloud.tf file written successfully."


# AWS云:
# -----------------------------------
# 创建AWS云opentofu模版文件夹
mkdir -p aws
# 进入
cd aws
# huaweicloud.tf文件并写入内容
cat << EOF > hashicorp.tf
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "4.65.0"
    }
  }
}

provider "aws" {
  # Configuration options
}
EOF
# 退出
cd ..
echo "Folder created and hashicorp.tf file written successfully."




# 执行tofu init 生成依赖项文件:
# -----------------------------------
# 初始化aliyun opentofu配置文件
cd /app/tofu-emplate/aliyun/
TF_PLUGIN_CACHE_DIR=/app/tofu-plugin tofu init

# 初始化tencent opentofu配置文件
cd /app/tofu-emplate/tencent/
TF_PLUGIN_CACHE_DIR=/app/tofu-plugin tofu init

# 初始化baidu opentofu配置文件
cd /app/tofu-emplate/baidu/
TF_PLUGIN_CACHE_DIR=/app/tofu-plugin tofu init

# 初始化huawei opentofu配置文件
cd /app/tofu-emplate/huawei/
TF_PLUGIN_CACHE_DIR=/app/tofu-plugin tofu init

# 初始化aws opentofu配置文件
cd /app/tofu-emplate/aws/
TF_PLUGIN_CACHE_DIR=/app/tofu-plugin tofu init