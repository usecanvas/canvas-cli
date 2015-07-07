require 'minitest/autorun'
require_relative 'helper'

class PullTest < MiniTest::Unit::TestCase
  HTML  = '<h1 id="hello-world-">Hello World!</h1>'
  MD    = '# Hello World!'
  CJSON = '{"type":"canvas","meta":{"tags":[]},"content":[{"type":"heading","content":"Hello World!","meta":{"level":1}}]}'

  def setup
		super
    @c_url = `echo "#{MD}" | #{CLI.bin} new`.strip
    @c_id  = @c_url.split('/').last
  end

	def teardown
		super
		CLI.delete(@c_id)
	end

  def test_pull_with_no_format
    body = `#{CLI.bin} pull #{@c_id}`.strip
    assert_equal(MD, body)
  end

  def test_pull_with_html_format
    body = `#{CLI.bin} pull #{@c_id} --html`.strip
    assert_equal(HTML, body)
  end

  def test_pull_with_md_format
    body = `#{CLI.bin} pull #{@c_id} --md`.strip
    assert_equal(MD, body)
  end

  def test_pull_with_json_format
    body = `#{CLI.bin} pull #{@c_id} --json`.strip
    assert_equal(CJSON, body)
  end
end
