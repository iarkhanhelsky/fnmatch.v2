require 'test/unit'
require 'yaml'

class FnFixture
  attr_reader :want
  
  def initialize(hash)
    @pattern = hash['pattern']
    @input = hash['input']
    @flags = convert_flags(hash['flags'])
    @want = hash['want']
    @hash = hash
  end

  def to_s
    flags = @hash['flags'] || []
    if flags.size > 0 
      flags = flags.join(' | ').gsub('fnmatch.', 'File::')
    else
      flags = '0'
    end

    "fnmatch('#{@pattern}', '#{@input}', #{flags}) => #{@want}"
  end
  
  def convert_flags(list)
    return 0 if list.nil?
    
    flags = 0
    list.each do |f|
      case f
      when 'fnmatch.FNM_NOESCAPE'
        flags |= File::FNM_NOESCAPE
      when 'fnmatch.FNM_PATHNAME'
          flags |= File::FNM_PATHNAME
      when 'fnmatch.FNM_PERIOD', 'fnmatch.FNM_DOTMATCH'
          flags |= File::FNM_DOTMATCH
      when 'fnmatch.FNM_CASEFOLD'
          flags |= File::FNM_CASEFOLD
      when 'fnmatch.FNM_IGNORECASE'
          flags |= File::FNM_IGNORECASE
      when 'fnmatch.FNM_FILE_NAME'
          flags |= File::FNM_FILE_NAME
      when 'fnmatch.FNM_EXTGLOB'
          flags |= File::FNM_EXTGLOB
      else
        fail f
      end
    end
    
    flags
  end
  
  def fnmatch
    File.fnmatch(@pattern, @input, @flags)
  end
end

class TestFnmatch < Test::Unit::TestCase
  Dir[File.join(__dir__, "testdata/*/*.yaml")].each do |file|
    fixtures = YAML.load(IO.read(file)).map { |c| FnFixture.new(c) }
    suite = File.basename(File.dirname(file))
    test_name = "test_#{suite}_#{File.basename(file).tr('-.', '_')}"

    fixtures.each.with_index do |f, i|
      define_method "#{test_name}_#{i}".to_sym do
        desc = "fnmatch('#{}')"
        assert_equal f.want, f.fnmatch, f.to_s
      end
    end
  end
end