require 'minitest/autorun'
require_relative 'helper'

class DeleteTest < Minitest::Unit::TestCase

	def test_deletes_canvas
		id = CLI.new_canvas
    output = `#{CLI.bin} list`.strip
    assert_includes(output, id)
		`#{CLI.bin} delete #{id}`
		output = `#{CLI.bin} list #{id}`
    refute_includes(output, id)
	end
end
