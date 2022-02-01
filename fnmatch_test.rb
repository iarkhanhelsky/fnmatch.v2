require 'test/unit'
require 'yaml'

class FnFixture
  attr_reader :want
  
  def initialize(hash)
    @pattern = hash['pattern']
    @input = hash['input']
    @flags = convert_flags(hash['flags'])
    @want = hash['want']
  end
  
  def convert_flags(list)
    return 0 if list.nil?
    
    flags = 0
    list.each do |f|
    end
    
    flags
  end
  
  def fnmatch
    File.fnmatch(@pattern, @input, @flags)
  end
end

class TestFnmatch < Test::Unit::TestCase
  Dir["testdata/bsd/*.yaml"].each do |file|

    fixtures = YAML.load(IO.read(file)).map { |c| FnFixture.new(c) }
    test_name = 'test_' + File.basename(file).tr('-.', '_')

    fixtures.each.with_index do |f, i|
      define_method "#{test_name}_#{i}".to_sym do
        assert_equal f.want, f.fnmatch
      end
    end
  end
end