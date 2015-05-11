require 'minitest/autorun'

class EnvTest < MiniTest::Unit::TestCase
  def test_env
    `canvas env`.split("\n")
  end
end
