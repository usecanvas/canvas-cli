require 'minitest/autorun'
require_relative 'helper'

class EnvTest < MiniTest::Unit::TestCase
  def test_env
    `#{CLI.bin} env`.split("\n")
  end
end
