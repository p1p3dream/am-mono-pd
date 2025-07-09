# CloudFront origin access control for S3.
resource "aws_cloudfront_origin_access_control" "www_oac" {
  name                              = "www-oac"
  description                       = "OAC for website S3 bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# Cache policy for HTML files (short cache).
resource "aws_cloudfront_cache_policy" "html_policy" {
  name        = "html-policy-${var.module_id}"
  comment     = "Cache policy for HTML files with short TTL"
  default_ttl = 0
  max_ttl     = 60
  min_ttl     = 0

  parameters_in_cache_key_and_forwarded_to_origin {
    cookies_config {
      cookie_behavior = "none"
    }

    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "none"
    }

    enable_accept_encoding_brotli = true
    enable_accept_encoding_gzip   = true
  }
}

# Cache policy for JS/CSS assets with content hash (long cache).
resource "aws_cloudfront_cache_policy" "assets_policy" {
  name        = "assets-policy-${var.module_id}"
  comment     = "Cache policy for versioned assets with long TTL"
  default_ttl = 31536000 # 1 year.
  max_ttl     = 31536000 # 1 year.
  min_ttl     = 31536000 # 1 year.

  parameters_in_cache_key_and_forwarded_to_origin {
    cookies_config {
      cookie_behavior = "none"
    }

    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "none"
    }

    enable_accept_encoding_brotli = true
    enable_accept_encoding_gzip   = true
  }
}

# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution.
resource "aws_cloudfront_distribution" "main" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  aliases = [
    var.domains.abodemine_saas_whitelabel.name,
  ]

  # https://docs.aws.amazon.com/cdk/api/v2/python/aws_cdk.aws_cloudfront/PriceClass.html.
  price_class = "PriceClass_100"

  origin {
    domain_name              = aws_s3_bucket.www_bucket.bucket_regional_domain_name
    origin_id                = "S3-${aws_s3_bucket.www_bucket.id}"
    origin_access_control_id = aws_cloudfront_origin_access_control.www_oac.id
  }

  # ALB origin for API requests
  origin {
    domain_name = var.load_balancers.main.dns_name
    origin_id   = "ALB-origin"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  # Default behavior - route to S3 for static website with short cache for HTML
  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-${aws_s3_bucket.www_bucket.id}"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    cache_policy_id        = aws_cloudfront_cache_policy.html_policy.id
  }

  # Ordered cache behavior for JavaScript, CSS and other versioned assets (long cache).
  # ordered_cache_behavior {
  #   path_pattern             = "*.js"
  #   allowed_methods          = ["GET", "HEAD", "OPTIONS"]
  #   cached_methods           = ["GET", "HEAD"]
  #   target_origin_id         = "S3-${aws_s3_bucket.www_bucket.id}"
  #   viewer_protocol_policy   = "redirect-to-https"
  #   compress                 = true
  #   cache_policy_id          = aws_cloudfront_cache_policy.assets_policy.id
  # }

  # ordered_cache_behavior {
  #   path_pattern             = "*.css"
  #   allowed_methods          = ["GET", "HEAD", "OPTIONS"]
  #   cached_methods           = ["GET", "HEAD"]
  #   target_origin_id         = "S3-${aws_s3_bucket.www_bucket.id}"
  #   viewer_protocol_policy   = "redirect-to-https"
  #   compress                 = true
  #   cache_policy_id          = aws_cloudfront_cache_policy.assets_policy.id
  # }

  ordered_cache_behavior {
    path_pattern           = "assets/*"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-${aws_s3_bucket.www_bucket.id}"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    cache_policy_id        = aws_cloudfront_cache_policy.assets_policy.id
  }

  # Route /api/* requests to the ALB
  ordered_cache_behavior {
    path_pattern     = "/api/*"
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "ALB-origin"

    forwarded_values {
      query_string = true
      headers      = ["Host", "Origin", "Authorization"]

      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 0
    max_ttl                = 0
  }

  # SPA support - return index.html for 404s (client-side routing).
  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  # SSL certificate.
  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.abodemine_main.arn
    ssl_support_method  = "sni-only"
  }
}
