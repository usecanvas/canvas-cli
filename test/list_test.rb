require 'minitest/autorun'
require_relative 'helper'

class ListTest < MiniTest::Unit::TestCase
  def test_new_canvas_appears_in_list
    id = CLI.new_canvas
    output = `#{CLI.bin} list`.strip
    assert_includes(output, id)
  end
end
