package httputil2

import (
    "testing"
)

func TestDetectCharset(t *testing.T) {
    LogDebug = true
    data := []byte(`<!doctype html>
<html>
<head>
    <meta charset="gbk" />

<script>
    (function(w,d){
    try{
`)
    if charset := DetectCharset(data); charset != "gbk" {
        t.Error("charset: got", charset, "expected:", "gbk")
    }

    data = []byte(`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
  <head>
    <link rel="stylesheet" type="text/css" href="/css/style.css">
    <link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
    <meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1">

    <title>how to use <meta http-equiv="Content-Type" content="text"> | instructions how use Content-Type - meta tags seo search engines</title>
    <meta name="keywords" content="content, content-type, webpage, source, html, http-equiv, meta tag, html tag">
    <meta name="description" content="The html meta tag HTTP-EQUIV CONTENT-TYPE allows you to specify the media type (i.e. text/html) and the character set.">

    <!-- Metatags Copyright 1999-2015 | MetaTags.info sinds 1999 -->

    <meta name="revisit-after" content="7 days">
    <meta name="copyright" content="The Metatags Company Inc. - Miami">
    <meta name="author" content="the Metatags Company Inc. - seo services">
    <meta name="web_copy_date" content="26-09-2015">
    <meta name="expires" content="02-10-2027">
    <meta name="web_content_type" content="Tips">
    <meta name="author" content="Metatags.info The Metatags Company Inc.">
    <meta name="country" content="USA">
    <meta name="web_author" content="Meta Tags Editorial Department ">

    <meta name="reply-to" content="submit.searchengines.submit@metatags.info">
    <meta name="robots" content="index, follow">
    <meta name="resource-type" content="document">
    <meta name="classification" content="internet">
    <meta name="distribution" content="global">
    <meta name="rating" content="safe for kids">
    <meta name="doc-type" content="public">
    <meta name="Identifier-URL" content="http://www.metatags.info/">
    <meta name="subject" content="meta tags, analyzer results - analyze your website">

    <meta name="contactName" content="Meta tags">
    <meta name="contactOrganization" content="Metatags - The Metatags Company BV">
    <meta name="contactStreetAddress1" content="Hardwareweg 4">
    <meta name="contactZipcode" content="3821 BM">
    <meta nam`)
    if charset := DetectCharset(data); charset != "iso-8859-1" {
        t.Error("charset: got", charset, "expected ", "iso-8859-1")
    }
}

func TestExtractCharsetFromContentType(t *testing.T) {
    if charset, _, _ := CharsetFromContentType("text/html;charset=GBK"); charset != "gbk" {
        t.Error("charset: got", charset, "expected:", "gbk")
    }
}
