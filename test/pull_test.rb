require 'minitest/autorun'

class PullTest < MiniTest::Unit::TestCase
  HTML  = '<h1 id="hello-world-">Hello World!</h1>'
  MD    = '# Hello World!'
  CJSON = '[{"type":"heading","depth":1,"text":"Hello World!"}]'

  def setup
    @c_url = `echo "#{MD}" | canvas new`.strip
    @c_id  = @c_url.split('/').last
  end

  def test_pull_with_no_format
    body = `canvas pull #{@c_id}`.strip
    assert_equal(MD, body)
  end

  def test_pull_with_html_format
    body = `canvas pull #{@c_id} --format=html`.strip
    assert_equal(HTML, body)
  end

  def test_pull_with_md_format
    body = `canvas pull #{@c_id} --format=md`.strip
    assert_equal(MD, body)
  end

  def test_pull_with_json_format
    body = `canvas pull #{@c_id} --format=json`.strip
    assert_equal(CJSON, body)
  end
end
