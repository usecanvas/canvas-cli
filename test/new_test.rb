require 'minitest/autorun'
require_relative 'helper'

class NewTest < MiniTest::Unit::TestCase
  DOC = "# Foo"

	def teardown
		super
		CLI.delete(@id)
	end

  def test_new_with_no_args
    @id = CLI.new_canvas("#{CLI.bin} new")
    assert_equal('', CLI.pull_canvas(@id))
  end

  def test_new_from_STDIN
    @id = CLI.new_canvas("echo \"#{DOC}\" | #{CLI.bin} new")
    assert_equal(DOC, CLI.pull_canvas(@id))
  end

  def test_new_from_FILE
    @id = CLI.new_canvas("#{CLI.bin} new README.md")
    readme = File.read('README.md').strip
    assert_equal(readme, CLI.pull_canvas(@id))
  end
end
